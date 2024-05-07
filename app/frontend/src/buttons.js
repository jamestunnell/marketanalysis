import van from "vanjs-core"

const { button, div } = van.tags

const ButtonAct = ({text, onclick}) => {
    return div(
        button(
            {
                class: "px-4 bg-indigo-500 p-3 ml-3 rounded-lg text-white hover:bg-indigo-400",
                onclick: onclick,
            },
            text,
        )
    )
}

const ButtonCancel = ({text, onclick}) => {
    return div(
        button(
            {
                class: "px-4 bg-gray-100 p-3 rounded-lg text-black hover:bg-gray-200",
                onclick: onclick,
            },
            text,
        )
    )
}

export { ButtonAct, ButtonCancel };