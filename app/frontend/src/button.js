import van from "vanjs-core"

const { button, div } = van.tags

const Button = ({text, onclick}) => {
    return div(
        button(
            {
                class: "bg-transparent hover:bg-blue-500 text-blue-700 font-semibold hover:text-white py-2 px-4 border border-blue-500 hover:border-transparent rounded",
                onclick: onclick,
            },
            text,
        )
    )
}

export default Button