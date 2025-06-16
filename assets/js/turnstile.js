(function() {
	let turnstileToken = null;
	let turnstileWidgetId = null;
	let turnstileRendered = false;
	let formInteracted = false;

	// Initialize when DOM is ready
	function init() {
		const form = document.getElementById('newsletterForm');
		if (!form) return;

		setupFormInteractionListeners();
		setupFormSubmission();

		// If using HTMX, reinitialize after swaps
		if (window.htmx) {
			document.body.addEventListener('htmx:afterSwap', function(event) {
				if (event.detail.target.id === 'newsletter-subscription') {
					// Reset state
					turnstileToken = null;
					turnstileWidgetId = null;
					turnstileRendered = false;
					formInteracted = false;

					// Reinitialize
					setTimeout(init, 100);
				}
			});
		}
	}

	// Setup interaction listeners
	function setupFormInteractionListeners() {
		const emailInput = document.getElementById('newsletter-email');
		if (!emailInput) return;

		// Show Turnstile on interaction
		['focus', 'click', 'input'].forEach(eventType => {
			emailInput.addEventListener(eventType, handleFormInteraction, { once: true });
		});
	}

	// Handle form interaction
	function handleFormInteraction() {
		if (!formInteracted && window.turnstile) {
			formInteracted = true;
			renderTurnstile();
		}
	}

	// Render Turnstile widget
	function renderTurnstile() {
		if (turnstileRendered) return;

		const container = document.getElementById('turnstile-container');
		if (!container) return;

		const siteKey = container.dataset.sitekey;

		// Show container with animation
		container.classList.remove('hidden');
		container.style.opacity = '0';
		container.style.transform = 'translateY(-10px)';

		setTimeout(() => {
			container.style.transition = 'opacity 0.3s ease, transform 0.3s ease';
			container.style.opacity = '1';
			container.style.transform = 'translateY(0)';
		}, 10);

		turnstileWidgetId = window.turnstile.render('#turnstile-container', {
			sitekey: siteKey,
			callback: function(token) {
				turnstileToken = token;
				// Store token in hidden field
				const tokenField = document.getElementById('turnstile-token');
				if (tokenField) {
					tokenField.value = token;
				}
				enableSubmitButton();
			},
			'expired-callback': function() {
				turnstileToken = null;
				const tokenField = document.getElementById('turnstile-token');
				if (tokenField) {
					tokenField.value = '';
				}
				disableSubmitButton();
			},
			'error-callback': function() {
				console.error('Turnstile error');
				disableSubmitButton();
			},
			theme: 'auto', // Will respect your site's theme
			size: 'normal'
		});

		turnstileRendered = true;
		disableSubmitButton();
	}

	// Enable submit button
	function enableSubmitButton() {
		const btn = document.getElementById('newsletter-submit-btn');
		if (btn) {
			btn.disabled = false;
		}
	}

	// Disable submit button
	function disableSubmitButton() {
		const btn = document.getElementById('newsletter-submit-btn');
		if (btn) {
			btn.disabled = true;
		}
	}

	// Setup form submission
	function setupFormSubmission() {
		const form = document.getElementById('newsletterForm');
		if (!form) return;

		// For HTMX, we need to intercept the submission
		form.addEventListener('htmx:beforeRequest', function(event) {
			if (!turnstileToken && formInteracted) {
				event.preventDefault();
				event.detail.xhr.abort();

				// Show error message
				showError('Please complete the security verification.');
				return false;
			}

			// Show loading state
			const btn = document.getElementById('newsletter-submit-btn');
			if (btn) {
				const textSpan = btn.querySelector('.newsletter-btn-text');
				const spinnerSpan = btn.querySelector('.newsletter-btn-spinner');

				if (textSpan) textSpan.classList.add('hidden');
				if (spinnerSpan) spinnerSpan.classList.remove('hidden');
				btn.disabled = true;
			}
		});

		// Reset button state after response
		form.addEventListener('htmx:afterRequest', function(event) {
			const btn = document.getElementById('newsletter-submit-btn');
			if (btn) {
				const textSpan = btn.querySelector('.newsletter-btn-text');
				const spinnerSpan = btn.querySelector('.newsletter-btn-spinner');

				if (textSpan) textSpan.classList.remove('hidden');
				if (spinnerSpan) spinnerSpan.classList.add('hidden');
				btn.disabled = !turnstileToken;
			}
		});

		// For non-HTMX fallback
		form.addEventListener('submit', function(event) {
			if (!window.htmx) {
				if (!turnstileToken && formInteracted) {
					event.preventDefault();
					showError('Please complete the security verification.');
					return false;
				}
			}
		});
	}

	// Show error message
	function showError(message) {
		// You can implement a toast notification or inline error
		const form = document.getElementById('newsletterForm');
		if (!form) return;

		// Check if error div exists
		let errorDiv = form.querySelector('.turnstile-error');
		if (!errorDiv) {
			errorDiv = document.createElement('div');
			errorDiv.className = 'turnstile-error text-error text-sm mt-2';
			const container = document.getElementById('turnstile-container');
			if (container && container.parentNode) {
				container.parentNode.insertBefore(errorDiv, container.nextSibling);
			}
		}

		errorDiv.textContent = message;

		// Remove after 5 seconds
		setTimeout(() => {
			errorDiv.remove();
		}, 5000);
	}

	// Turnstile callback
	window.onloadTurnstileCallback = function() {
		// Turnstile is loaded, initialize if form is already present
		init();
	};

	// Initialize when DOM is ready
	if (document.readyState === 'loading') {
		document.addEventListener('DOMContentLoaded', init);
	} else {
		init();
	}
})();
