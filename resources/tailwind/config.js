/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "../views/internal/pages/**/*.templ",
    "../views/internal/layouts/**/*.templ",
    "../views/internal/components/**/*.templ",
    "node_modules/preline/dist/*.js",
  ],
  plugins: [
    require('@tailwindcss/forms'),
    require("@tailwindcss/typography"),
    require('preline/plugin')
  ],
}
