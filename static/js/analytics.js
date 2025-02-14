(function() {
  const script = document.currentScript;
  const websiteId = script.getAttribute('data-website-id');

  // Session and visitor tracking
  const getOrCreateId = (key) => {
    let id = localStorage.getItem(key);
    if (!id) {
      id = typeof crypto !== 'undefined' && crypto.randomUUID 
        ? crypto.randomUUID() 
        : Math.random().toString(36).substring(2, 15);
      localStorage.setItem(key, id);
    }
    return id;
  };

  // Bot/crawler detection
  const isBotOrCrawler = () => {
    const botPatterns = [
      /bot/i, /crawl/i, /spider/i, 
      /phantom/i, /nightmare/i, 
      /webdriver/i, /cypress/i
    ];
    return botPatterns.some(pattern => 
      pattern.test(navigator.userAgent) || 
      window.navigator.webdriver
    );
  };

  // Scroll tracking
  let maxScrollDepth = 0;
  const updateScrollDepth = () => {
    const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
    const windowHeight = window.innerHeight;
    const documentHeight = Math.max(
      document.body.scrollHeight, 
      document.body.offsetHeight, 
      document.documentElement.clientHeight, 
      document.documentElement.scrollHeight, 
      document.documentElement.offsetHeight
    );
    const scrollPercent = Math.round((scrollTop / (documentHeight - windowHeight)) * 100);
    maxScrollDepth = Math.max(maxScrollDepth, scrollPercent);
  };

  // Page exit tracking
  const trackPageExit = () => {
    const payload = {
	  website_id: websiteId,
      type: 'pageleave',
      url: window.location.href,
      scroll_depth: maxScrollDepth,
      title: document.title
    };
    sendEvent(payload, 'pageleave');
  };

  // Send event function
  function sendEvent(payload, eventType = 'event') {
    // Skip for bots and local development
    if (isBotOrCrawler() || 
        /^localhost$|^127\./.test(window.location.hostname) || 
        window.localStorage.getItem('analytics-disabled')) {
      return;
    }

    const finalPayload = {
	  website_id: websiteId,
      type: eventType,
      url: window.location.href,
      path: window.location.pathname,
      referrer: document.referrer,
      title: document.title,
      timestamp: new Date().toISOString(),
      screen: `${window.screen.width}x${window.screen.height}`,
      language: navigator.language || navigator.userLanguage,
      visitor_id: getOrCreateId('analytics-visitor-id'),
      session_id: getOrCreateId('analytics-session-id'),
      ...payload
    };

    try {
      navigator.sendBeacon('/api/v1/collect', JSON.stringify(finalPayload));
    } catch {
      fetch('/api/v1/collect', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(finalPayload)
      });
    }
  }

  // Initialize tracking
  if (websiteId) {
    // Track page view
    sendEvent({}, 'pageview');

    // Track clicks
    document.addEventListener('click', (event) => {
      const target = event.target.closest('a, button, [data-analytics]');
      if (target) {
        const payload = {
          element_tag: target.tagName.toLowerCase(),
          element_text: target.textContent?.trim(),
          element_class: target.className,
          element_id: target.id,
          element_href: target.href,
          custom_data: target.getAttribute('data-analytics')
        };
        sendEvent(payload, 'click');
      }
    });

    // Scroll tracking
    window.addEventListener('scroll', updateScrollDepth);

    // Page exit tracking
    window.addEventListener('beforeunload', trackPageExit);

    // Handle history changes for SPA
    window.addEventListener('popstate', () => {
      sendEvent({}, 'pageview');
    });
  }
})();
