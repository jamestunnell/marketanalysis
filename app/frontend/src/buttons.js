import van from "vanjs-core"
import { Tooltip } from 'vanjs-ui'
import { Toggle } from "vanjs-ui"
import { validateConnection } from "./connection"

const { button, div, input, label, span } = van.tags

const toStyleStr = (style) => Object.entries(style).map(([k, v]) => `${k}: ${v};`).join("");

const Button = ({child, onclick, disabled}) => {
    return buttonDisableable({
        child: child,
        onclick: onclick,
        disabled: disabled,
        classNormal: 'rounded-md p-1 m-1 text-white bg-gray-500 hover:bg-gray-600',
        classDisabled: 'rounded-md p-1 m-1 text-white bg-gray-300 cursor-not-allowed',
    })
}

const ButtonDanger = ({child, onclick, disabled}) => {
    return buttonDisableable({
        child: child,
        onclick: onclick,
        disabled: disabled,
        classNormal: 'rounded-md p-1 m-1 text-white bg-pink-500 hover:bg-pink-600',
        classDisabled: 'rounded-md p-1 m-1 text-white bg-pink-300 cursor-not-allowed',
    })
}

const iconClassNormal = 'rounded-md p-1 m-1 text-gray-700 hover:text-black'

const ButtonIconDisableable = ({icon, text, disabled, onclick}) => {
    return button(
        {
            class: van.derive(() => {
                return disabled.val ? 'rounded-md p-1 m-1 text-gray-400' : iconClassNormal
            }),
            disabled: disabled,
            onclick: onclick,
        },
        icon,
        text,
    )
}

const ButtonIcon = ({icon, text, onclick}) => {
    return button(
        {
            class: iconClassNormal,
            onclick: onclick,
        },
        icon,
        text,
    )
}

const toggleInputStyle = toStyleStr({
    opacity: 0,
    width: 0,
    height: 0,
    position: "absolute",
    "z-index": 10000, // Ensures the toggle clickable
})

const ButtonToggle = ({setVal, onSet, onClear}) => {
    const set = van.state(setVal)
    const offColor = "#9CA3AF"  // equivalent to text-gray-400 or rgb(156, 163, 175)
    const onColor = "#374151"   // equivalent to text-gray-700 or rgb(55, 65, 81)
    const size = 1
    const toggleLabelStyles = toStyleStr({
        position: "relative",
        display: "inline-block",
        width: 1.76 * size + "rem",
        height: size + "rem",
        cursor: "pointer",
    })
    const toggleSliderStyles = toStyleStr({
        position: "absolute",
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        transition: ".4s",
        "border-radius": size + "rem",
    })
    const toggleCircleStyles = toStyleStr({
        position: "absolute",
        height: 0.76 * size + "rem",
        width: 0.76 * size + "rem",
        left: 0.12 * size + "rem",
        bottom: 0.12 * size + "rem",
        "background-color": "white",
        transition: ".4s",
        "border-radius": "50%",
    });
    const toggleCircleTransl = toStyleStr({
        transform: `translateX(${0.76 * size}rem)`,
    })

    return div(
        {class: "rounded-md p-1 mt-3 mb-3"},
        label(
            {style: toggleLabelStyles},
            input({
                type: "checkbox",
                style: toggleInputStyle,
                oninput: e => {
                    set.val = e.target.checked

                    if (set.val) {
                        onSet()
                    } else {
                        onClear()
                    }
                },
            }), 
            span(
                {style: () => `${toggleSliderStyles}; background-color: ${set.val ? onColor : offColor};`},
                span({ style: () => toggleCircleStyles + (set.val ? toggleCircleTransl : "")}),
            ),
        )
    )
}

const ButtonIconTooltip = ({icon, text, tooltipText}) => {
    const showTooltip = van.state(false)
    
    return button(
        {
            style: "position: relative;",
            class: iconClassNormal,
            onmouseenter: () => showTooltip.val = true,
            onmouseleave: () => showTooltip.val = false,
        },
        icon,
        text,
        Tooltip({
            text: tooltipText,
            show: showTooltip,
        }),
    )
}

const ButtonCancel = ({child, onclick}) => {
    return button(
        {
            class: 'rounded-md p-1 m-3 text-black bg-gray-200 hover:bg-gray-300',
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

export { Button, ButtonDanger, ButtonCancel, ButtonIcon, ButtonIconDisableable, ButtonIconTooltip, ButtonToggle };