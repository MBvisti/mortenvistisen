import "./easymde_2-26-0.min.js";

(function() {
	const targetElement = document.getElementById('editorTarget');
	if (!targetElement) return;

	const easyMDE = new EasyMDE({
		element: targetElement,
		maxHeight: '750px',
		toolbar: [
			"bold", "italic", "strikethrough", "|",
			"heading-2", "heading-3", "|",
			"quote", "unordered-list", "ordered-list", "|",
			"link", "image", "code", "|",
			"preview", "fullscreen",
		],
		indentWithTabs: true,
		tabSize: 4,
	});

	// Sync editor content back to textarea before form submission
	const form = targetElement.closest('form');
	if (form) {
		form.addEventListener('submit', function() {
			targetElement.value = easyMDE.value();
		});
	}

	// Also sync on editor change for data-bind compatibility
	easyMDE.codemirror.on('change', function() {
		targetElement.value = easyMDE.value();
		// Dispatch input event for any data-binding frameworks
		targetElement.dispatchEvent(new Event('input', { bubbles: true }));
	});
})()
