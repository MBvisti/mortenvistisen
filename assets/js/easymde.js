import "./easymde-2_26_0.min.js";

(function() {
	const element = document.getElementById('my-text-area');
	if (element) {
		new window.EasyMDE({
			element: element,
			placeholder: "Write your article content in Markdown...",
			spellChecker: false,
			hideIcons: ["guide", "fullscreen", "side-by-side"],
			showIcons: ["code", "table"],
			status: ["autosave", "lines", "words", "cursor"],
			tabSize: 4,
			indentWithTabs: false,
			lineWrapping: true,
			autofocus: false,
			toolbar: [
				"bold", "italic", "strikethrough", "|",
				"heading-1", "heading-2", "heading-3", "|",
				"quote", "unordered-list", "ordered-list", "|",
				"link", "image", "code", "table", "|",
				"preview", "guide"
			]
		});
	}
})()

