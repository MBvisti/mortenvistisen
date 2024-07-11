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
  daisyui: {
    themes: ["dark"],
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
    require('daisyui')
  ],
}
