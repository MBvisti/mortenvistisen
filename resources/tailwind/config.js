/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "../views/**/*.templ",
    "../views/***/**/*.templ",
    "../views/internal/layouts/**/*.templ",
    "../views/internal/components/**/*.templ",
    "../posts/**/*.md"
  ],
  darkMode: 'class',
  daisyui: {
    themes: ["dracula"],
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
    require('daisyui')
  ],
}
