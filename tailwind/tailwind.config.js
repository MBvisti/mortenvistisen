module.exports = {
    darkMode: "class",
    content: ["../templates/**/*.{html,js}"],
    daisyui: {
        theme: ['garden', 'dracula']
    },
    theme: {
        extend: {
            fontFamily: {
                'sans': ['Montserrat', 'Helvetica', 'Arial', 'sans-serif']
            },
        }
    },
    plugins: [require("@tailwindcss/typography") ,require('daisyui')],
}
