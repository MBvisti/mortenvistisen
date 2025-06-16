// Newsletter Subscription JavaScript - Fixed Version
// Handles Turnstile CAPTCHA integration and form enhancement

class NewsletterSubscription {
	constructor() {
		this.turnstileLoaded = false;
		this.turnstileWidget = null;
		this.scriptLoadInProgress = false;
		this.init();
	}

	init() {
		document.addEventListener('DOMContentLoaded', () => {
			this.setupEventListeners();
			this.enhanceForm();
		});
	}

	setupEventListeners() {
		// Handle email input focus to load Turnstile
		document.addEventListener('focus', (e) => {
			if (e.target.id === 'newsletter-email') {
				this.loadTurnstile();
			}
		}, true);

		// Handle form submission
		document.addEventListener('htmx:beforeRequest', (e) => {
			if (e.target.getAttribute('hx-post') === '/subscribe') {
				this.handleFormSubmission(e);
			}
		});

		// Handle HTMX responses
		document.addEventListener('htmx:afterRequest', (e) => {
			if (e.target.getAttribute('hx-post') === '/subscribe') {
				this.handleFormResponse(e);
			}
		});
	}

	enhanceForm() {
		const forms = document.querySelectorAll('form[hx-post="/subscribe"]');
		forms.forEach(form => {
			// Add loading state handling
			const submitBtn = form.querySelector('.newsletter-submit-btn');
			if (submitBtn) {
				submitBtn.addEventListener('click', () => {
					this.setLoadingState(submitBtn, true);
				});
			}
		});
	}

	loadTurnstile() {
		if (this.turnstileLoaded || this.scriptLoadInProgress) return;

		// Check if script is already in DOM (with or without query params)
		const existingScript = document.querySelector('script[src*="challenges.cloudflare.com/turnstile/v0/api.js"]');
		if (existingScript) {
			this.turnstileLoaded = true;
			// Wait a bit for script to fully load
			setTimeout(() => this.initializeTurnstile(), 100);
			return;
		}

		this.scriptLoadInProgress = true;
		const script = document.createElement('script');
		script.src = 'https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit';
		script.async = true;
		script.defer = true;
		script.onload = () => {
			this.turnstileLoaded = true;
			this.scriptLoadInProgress = false;
			this.initializeTurnstile();
		};
		script.onerror = () => {
			this.scriptLoadInProgress = false;
			console.error('Failed to load Turnstile script');
		};
		document.head.appendChild(script);
	}

	initializeTurnstile() {
		// Wait for Turnstile to be available
		const checkTurnstile = () => {
			if (typeof window.turnstile !== 'undefined') {
				const captchaElements = document.querySelectorAll('.cf-turnstile:not([data-turnstile-initialized])');
				captchaElements.forEach(element => {
					try {
						// Mark as initialized immediately to prevent race conditions
						element.setAttribute('data-turnstile-initialized', 'true');

						const widgetId = window.turnstile.render(element, {
							sitekey: element.getAttribute('data-sitekey'),
							callback: (token) => {
								this.onTurnstileSuccess(token, element);
							},
							'error-callback': () => {
								this.onTurnstileError(element);
							},
							'expired-callback': () => {
								this.onTurnstileExpired(element);
							}
						});

						// Store widget ID on the element for future reference
						element.setAttribute('data-turnstile-widget-id', widgetId);

						console.log('Turnstile initialized successfully, widget ID:', widgetId);

					} catch (error) {
						console.error('Failed to initialize Turnstile:', error);
						// Remove the initialized flag if initialization failed
						element.removeAttribute('data-turnstile-initialized');
					}
				});
			} else {
				setTimeout(checkTurnstile, 100);
			}
		};
		checkTurnstile();
	}

	onTurnstileSuccess(token, element) {
		console.log('Turnstile success, token:', token ? 'received' : 'missing');

		// Store the token in a hidden input
		const form = element.closest('form[hx-post="/subscribe"]');
		if (form) {
			// Remove any existing token input
			const existingInput = form.querySelector('input[name="cf-turnstile-response"]');
			if (existingInput) {
				existingInput.remove();
			}

			// Create new hidden input with token
			const tokenInput = document.createElement('input');
			tokenInput.type = 'hidden';
			tokenInput.name = 'cf-turnstile-response';
			tokenInput.value = token;
			form.appendChild(tokenInput);

			// Enable submit button
			const submitBtn = form.querySelector('.newsletter-submit-btn');
			if (submitBtn) {
				submitBtn.disabled = false;
			}

			// Remove any validation errors
			const validationError = form.querySelector('.validation-error');
			if (validationError) {
				validationError.remove();
			}
		}
	}

	onTurnstileError(element) {
		console.error('Turnstile error occurred');
		this.showTurnstileError(element);
	}

	onTurnstileExpired(element) {
		console.warn('Turnstile token expired');
		// Remove the token input since it's no longer valid
		const form = element.closest('form[hx-post="/subscribe"]');
		if (form) {
			const tokenInput = form.querySelector('input[name="cf-turnstile-response"]');
			if (tokenInput) {
				tokenInput.remove();
			}

			// Disable form submission until new token is obtained
			const submitBtn = form.querySelector('.newsletter-submit-btn');
			if (submitBtn) {
				submitBtn.disabled = true;
			}
		}
	}

	showTurnstileError(element) {
		const form = element.closest('form[hx-post="/subscribe"]');
		if (form) {
			// Create error message if it doesn't exist
			let errorMsg = form.querySelector('.turnstile-error');
			if (!errorMsg) {
				errorMsg = document.createElement('div');
				errorMsg.className = 'turnstile-error bg-error text-error-content p-4 rounded-lg mb-4 text-center';
				errorMsg.innerHTML = '<p class="font-medium">There was an issue with the security verification. Please try again.</p>';
				element.insertAdjacentElement('beforebegin', errorMsg);
			}
		}
	}

	handleFormSubmission(e) {
		console.log('Form submission started');
		const form = e.target;
		const submitBtn = form.querySelector('.newsletter-submit-btn');

		// Check for the token in the hidden input
		const tokenInput = form.querySelector('input[name="cf-turnstile-response"]');
		console.log('Token input found:', !!tokenInput);
		console.log('Token value:', tokenInput ? (tokenInput.value ? 'present' : 'empty') : 'no input');

		const turnstileResponse = tokenInput ? tokenInput.value : null;

		if (!turnstileResponse) {
			e.preventDefault();
			this.showValidationError(form, 'Please complete the security verification.');
			if (submitBtn) {
				this.setLoadingState(submitBtn, false);
			}
			return false;
		}

		// Set loading state
		if (submitBtn) {
			this.setLoadingState(submitBtn, true);
		}

		return true;
	}

	handleFormResponse(e) {
		const form = e.target;
		const submitBtn = form.querySelector('.newsletter-submit-btn');

		// Reset loading state
		if (submitBtn) {
			this.setLoadingState(submitBtn, false);
		}

		// Check if there was an error (form is still present)
		if (form.closest('.newsletter-subscription')) {
			// Reset Turnstile for retry
			this.resetTurnstile(form);
		}
	}

	showValidationError(form, message) {
		// Remove existing error
		const existingError = form.querySelector('.validation-error');
		if (existingError) {
			existingError.remove();
		}

		// Create new error message
		const errorDiv = document.createElement('div');
		errorDiv.className = 'validation-error bg-error text-error-content p-4 rounded-lg mb-4 text-center';
		errorDiv.innerHTML = `<p class="font-medium">${message}</p>`;

		const firstChild = form.firstElementChild;
		form.insertBefore(errorDiv, firstChild);

		// Auto-remove after 5 seconds
		setTimeout(() => {
			if (errorDiv.parentNode) {
				errorDiv.remove();
			}
		}, 5000);
	}

	setLoadingState(button, loading) {
		if (!button) return;

		const textSpan = button.querySelector('.newsletter-btn-text');
		const loadingSpan = button.querySelector('.hidden');

		if (loading) {
			button.disabled = true;
			if (textSpan) textSpan.style.display = 'none';
			if (loadingSpan) {
				loadingSpan.classList.remove('hidden');
				loadingSpan.style.display = 'flex';
			}
		} else {
			button.disabled = false;
			if (textSpan) textSpan.style.display = 'inline';
			if (loadingSpan) {
				loadingSpan.classList.add('hidden');
				loadingSpan.style.display = 'none';
			}
		}
	}

	resetTurnstile(form) {
		if (typeof window.turnstile !== 'undefined') {
			const turnstileElement = form.querySelector('.cf-turnstile[data-turnstile-widget-id]');
			if (turnstileElement) {
				const widgetId = turnstileElement.getAttribute('data-turnstile-widget-id');
				try {
					window.turnstile.reset(widgetId);
					// Also remove any existing token input
					const tokenInput = form.querySelector('input[name="cf-turnstile-response"]');
					if (tokenInput) {
						tokenInput.remove();
					}
				} catch (error) {
					console.error('Failed to reset Turnstile:', error);
				}
			}
		}
	}

	// Static method to reset Turnstile globally (for templ script)
	static resetCaptcha() {
		if (typeof window.turnstile !== 'undefined') {
			try {
				// Reset all widgets on the page
				document.querySelectorAll('.cf-turnstile[data-turnstile-widget-id]').forEach(element => {
					const widgetId = element.getAttribute('data-turnstile-widget-id');
					if (widgetId) {
						window.turnstile.reset(widgetId);
					}
				});
			} catch (error) {
				console.error('Failed to reset Turnstile globally:', error);
			}
		}
	}
}

// Initialize the newsletter subscription handler
const newsletterSubscription = new NewsletterSubscription();

// Make resetCaptcha available globally for templ scripts
window.resetCaptcha = NewsletterSubscription.resetCaptcha;
