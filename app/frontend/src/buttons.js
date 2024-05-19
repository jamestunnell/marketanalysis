import van from "vanjs-core"
import { Tooltip } from 'vanjs-ui'

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

const iconClassNormal = 'rounded-md p-3 m-1 text-gray-700 hover:text-black'

const ButtonIconDisableable = ({icon, disabled, onclick}) => {
    return button(
        {
            class: van.derive(() => {
                return disabled.val ? 'rounded-md p-3 m-1 text-gray-400' : iconClassNormal
            }),
            disabled: disabled,
            onclick: onclick,
        },
        icon,
    )
}

const ButtonIcon = ({icon, onclick}) => {
    return button(
        {
            class: iconClassNormal,
            onclick: onclick,
        },
        icon,
    )
}

const ButtonIconTooltip = ({icon, tooltipText}) => {
    const showTooltip = van.state(false)
    
    return button(
        {
            style: "position: relative;",
            class: iconClassNormal,
            onmouseenter: () => showTooltip.val = true,
            onmouseleave: () => showTooltip.val = false,
        },
        icon,
        Tooltip({
            text: tooltipText,
            show: showTooltip,
        }),
    )
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
            onmouseenter: () => showTooltip.val = true,
            onmouseleave: () => showTooltip.val = false,
            onclick: onclick,
            disabled: disabled,
        },
        child,
    );
}

export { Button, ButtonDanger, ButtonCancel, ButtonIcon, ButtonIconDisableable, ButtonIconTooltip };