/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "../views/**/*.templ",
  ],
  darkMode: 'class',
  corePlugins: {
	container: false,
  },
  daisyui: {
    themes: ["dark"],
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
    require('daisyui'),
	require('tailwind-bootstrap-grid')({
      containerMaxWidths: {
        sm: '540px',
        md: '720px',
        lg: '960px',
        xl: '1140px',
      },
    }),
  ],
}
