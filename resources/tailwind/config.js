/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "../views/**/*.templ",
    "../views/internal/layouts/**/*.templ",
    "../views/internal/components/**/*.templ",
    "../posts/**/*.md",
    "node_modules/preline/dist/*.js",
  ],
  darkMode: 'class',
  plugins: [
    require('@tailwindcss/forms'),
    require("@tailwindcss/typography"),
    require('preline/plugin')
  ],
}
