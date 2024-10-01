/** @type {import('tailwindcss').Config} */
module.exports = {
	content: [
		"../views/*.templ",
		"../views/authentication/*.templ",
		"../views/components/*.templ",
		"../views/internal/components/*.templ",
		"../views/dashboard/*.templ",
		"../views/internal/layouts/*.templ"
	],
	darkMode: 'class',
	daisyui: {
		themes: ["halloween"],
	},
	plugins: [
		require('@tailwindcss/forms'),
		require('@tailwindcss/typography'),
		require('daisyui')
	],
}
