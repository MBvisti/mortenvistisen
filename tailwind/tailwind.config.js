module.exports = {
    content: ["../templates/**/*.{html,js}"],
    daisyui: {
        theme: ['dracular']
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
