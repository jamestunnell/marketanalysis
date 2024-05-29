import van from "vanjs-core"

const {div} = van.tags

const ModalBackground = (...children) => {
    return div(
        {
            class: "flex items-center justify-center fixed left-0 right-0 top-0 bottom-0 z-10",
            style: "background-color:rgba(0,0,0,.5);"
        },
        children,
    )
}

const ModalForeground = ({...props}, ...children) => {
    return div({class: "block p-16 rounded-lg bg-white", ...props}, children)
};

export {ModalBackground, ModalForeground};