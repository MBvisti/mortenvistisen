/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "../views/*.templ",
    "../views/authentication/*.templ",
    "../views/components/*.templ",
    "../views/internal/components/*.templ",
    "../views/internal/layouts/*.templ"
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
