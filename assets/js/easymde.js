import "./easymde_2-26-0.min.js";

(function() {
	new EasyMDE({
		element: document.getElementById('editorTarget'),
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
})()
