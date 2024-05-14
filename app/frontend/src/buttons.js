import van from "vanjs-core"

const { button } = van.tags

const Button = ({child, onclick, disabled}) => {
    return buttonDisableable({
        child: child,
        onclick: onclick,
        disabled: disabled,
        classNormal: 'rounded-md p-3 m-1 text-white bg-indigo-500 hover:bg-indigo-600',
        classDisabled: 'rounded-md p-3 m-1 text-white bg-indigo-300 cursor-not-allowed',
    })
}

const ButtonDanger = ({child, onclick, disabled}) => {
    return buttonDisableable({
        child: child,
        onclick: onclick,
        disabled: disabled,
        classNormal: 'rounded-md p-3 m-1 text-white bg-pink-500 hover:bg-pink-600',
        classDisabled: 'rounded-md p-3 m-1 text-white bg-pink-300 cursor-not-allowed',
    })
}

const ButtonCancel = ({child, onclick}) => {
    return button(
        {
            class: 'rounded-md p-3 m-1 text-black bg-gray-200 hover:bg-gray-300',
            onclick: onclick,
        },
        child,
    );
}

const buttonDisableable = ({child, onclick, disabled, classNormal, classDisabled}) => {
    return button(
        {
            class: van.derive(() => {
                return disabled ? (disabled.val ? classDisabled : classNormal) : classNormal
            }),
            onclick: onclick,
            disabled: disabled,
        },
        child,
    );
}

export { Button, ButtonDanger, ButtonCancel };