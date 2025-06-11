// Newsletter Subscription JavaScript
// Handles Turnstile CAPTCHA integration and form enhancement

class NewsletterSubscription {
    constructor() {
        this.turnstileLoaded = false;
        this.turnstileWidget = null;
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
        if (this.turnstileLoaded) return;

        const script = document.createElement('script');
        script.src = 'https://challenges.cloudflare.com/turnstile/v0/api.js';
        script.async = true;
        script.defer = true;
        script.onload = () => {
            this.turnstileLoaded = true;
            this.initializeTurnstile();
        };
        document.head.appendChild(script);
    }

    initializeTurnstile() {
        // Wait for Turnstile to be available
        const checkTurnstile = () => {
            if (typeof window.turnstile !== 'undefined') {
                const captchaElements = document.querySelectorAll('.cf-turnstile');
                captchaElements.forEach(element => {
                    if (!element.hasAttribute('data-turnstile-initialized')) {
                        try {
                            this.turnstileWidget = window.turnstile.render(element, {
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
                            element.setAttribute('data-turnstile-initialized', 'true');
                        } catch (error) {
                            console.error('Failed to initialize Turnstile:', error);
                        }
                    }
                });
            } else {
                setTimeout(checkTurnstile, 100);
            }
        };
        checkTurnstile();
    }

    onTurnstileSuccess(token, element) {
        // Enable form submission
        const form = element.closest('form[hx-post="/subscribe"]');
        if (form) {
            const submitBtn = form.querySelector('.newsletter-submit-btn');
            if (submitBtn) {
                submitBtn.disabled = false;
            }
        }
    }

    onTurnstileError(element) {
        console.error('Turnstile error occurred');
        this.showTurnstileError(element);
    }

    onTurnstileExpired(element) {
        console.warn('Turnstile token expired');
        // Disable form submission until new token is obtained
        const form = element.closest('form[hx-post="/subscribe"]');
        if (form) {
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
        const form = e.target;
        const submitBtn = form.querySelector('.newsletter-submit-btn');
        
        // Validate Turnstile
        const turnstileResponse = this.getTurnstileResponse(form);
        if (!turnstileResponse) {
            e.preventDefault();
            this.showValidationError(form, 'Please complete the security verification.');
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

        // Reset Turnstile if form is still present (error case)
        if (form.closest('.newsletter-subscription')) {
            this.resetTurnstile(form);
        }
    }

    getTurnstileResponse(form) {
        const turnstileElement = form.querySelector('.cf-turnstile');
        if (!turnstileElement) return null;

        // Get the Turnstile response token
        const iframe = turnstileElement.querySelector('iframe');
        if (!iframe) return null;

        // Check if Turnstile has provided a token
        const input = form.querySelector('input[name="cf-turnstile-response"]');
        return input ? input.value : null;
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
        if (typeof window.turnstile !== 'undefined' && this.turnstileWidget) {
            try {
                window.turnstile.reset(this.turnstileWidget);
            } catch (error) {
                console.error('Failed to reset Turnstile:', error);
            }
        }
    }

    // Static method to reset Turnstile globally (for templ script)
    static resetCaptcha() {
        if (typeof window.turnstile !== 'undefined') {
            try {
                window.turnstile.implicitRender();
            } catch (error) {
                console.error('Failed to reset Turnstile implicitly:', error);
            }
        }
    }
}

// Initialize the newsletter subscription handler
const newsletterSubscription = new NewsletterSubscription();

// Make resetCaptcha available globally for templ scripts
window.resetCaptcha = NewsletterSubscription.resetCaptcha;