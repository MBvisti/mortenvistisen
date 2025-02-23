(function() {
  const script = document.currentScript;

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

  const trackPageExit = () => {
    const payload = {
      type: 'page_leave',
      url: window.location.href,
      scroll_depth: maxScrollDepth,
      title: document.title
    };
    sendEvent(payload, 'page_leave');
  };

  function sendEvent(payload, eventType = 'event') {
    if (isBotOrCrawler() || 
        /^localhost$|^127\./.test(window.location.hostname) || 
        window.localStorage.getItem('analytics-disabled')) {
      return;
    }

    const finalPayload = {
      type: eventType,
      url: window.location.href,
      path: window.location.pathname,
	  user_agent: navigator.userAgent,
      referrer: document.referrer,
      title: document.title,
      timestamp: new Date().toISOString(),
      screen: `${window.screen.width}x${window.screen.height}`,
      language: navigator.language || navigator.userLanguage,
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

    sendEvent({}, 'page_view');

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

    window.addEventListener('scroll', updateScrollDepth);

    window.addEventListener('beforeunload', trackPageExit);

    window.addEventListener('popstate', () => {
      sendEvent({}, 'page_view');
    });
})();
