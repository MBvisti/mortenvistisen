module.exports = {
    content: ["../templates/**/*.{html,js}"],
    theme: {
        extend: {
            fontFamily: {
                'sans': ['Montserrat', 'Helvetica', 'Arial', 'sans-serif']
            },
        },
    },
    plugins: [require("@tailwindcss/typography"), require('daisyui')],
    daisyui: {
        theme: ['lofi']
    },
}
