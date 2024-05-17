import van from "vanjs-core"
import { nanoid } from 'nanoid'

const {button, div, h3} = van.tags

import { IconError } from "./icons"

const AlertError = ({title="Error", substance}) => {
    const alertID = `alert-${nanoid()}`;
    const buttonID = `dismiss-${nanoid()}`;
    const dismissBtn = button(
        {
            id: buttonID,
            type:"button",
            class: "text-red-800 bg-transparent border border-red-800 hover:bg-red-900 hover:text-white focus:ring-4 focus:outline-none focus:ring-red-300 font-medium rounded-lg text-xs px-3 py-1.5 text-center dark:hover:bg-red-600 dark:border-red-600 dark:text-red-500 dark:hover:text-white dark:focus:ring-red-800",
        },
        "Dismiss"
    )

    dismissBtn.setAttribute('data-dismiss-target', `#${alertID}`)
    dismissBtn.setAttribute('aria-label', 'Close')

    const autoDismiss = () => {
        const btn = document.getElementById(buttonID);
        if (btn) {
            btn.click()
        }
    }

    setTimeout(autoDismiss, 20000);

    return div(
        {
            id:alertID,
            class: "p-4 mb-4 text-red-800 border border-red-300 rounded-lg bg-red-50 dark:bg-gray-800 dark:text-red-400 dark:border-red-800",
            role: "alert",
        },
        div(
            {class:"flex items-center"},
            IconError(),
            h3({class:"text-lg font-medium"}, title),
        ),
        div(
            {class:"mt-2 mb-4 text-sm"},
            substance
        ),
        div(
            {class:"flex"},
            dismissBtn
        ),
    )
}

export {AlertError}