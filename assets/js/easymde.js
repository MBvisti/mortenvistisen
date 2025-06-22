import "./easymde-2_26_0.min.js";

(function() {
	new EasyMDE({
		element: document.getElementById('editorTarget'),
		maxHeight: '500px',
		toolbar: [
			"bold", "italic", "strikethrough", "|",
			"heading-2", "heading-3", "|",
			"quote", "unordered-list", "ordered-list", "|",
			"link", "image", "code", "|",
			"preview", "fullscreen",
		],
		indentWithTabs: true,
		// lineNumbers: true,
		tabSize: 4,
	});
	// const element = document.getElementById('my-text-area');
	// if (element) {
	// 	new window.EasyMDE({
	// 		element: element,
	// 		// placeholder: "Write your article content in Markdown...",
	// 		spellChecker: false,
	// 		initialEditType: "markdown",
	// 		initialValue: "",
	// 		lineNumbers: true,
	// 		hideIcons: ["guide", "fullscreen", "side-by-side"],
	// 		showIcons: ["code", "table"],
	// 		status: ["autosave", "lines", "words", "cursor"],
	// 		tabSize: 4,
	// 		indentWithTabs: true,
	// 		lineWrapping: true,
	// 		usageStatistics: false,
	// 		autofocus: false,
	// 		toolbar: [
	// 			"bold", "italic", "strikethrough", "|",
	// 			"heading-1", "heading-2", "|",
	// 			"quote", "unordered-list", "ordered-list", "|",
	// 			"link", "image", "code", "|",
	// 			"preview",
	// 		]
	// 	});
	// }
})()

