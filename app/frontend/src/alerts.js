import van from "vanjs-core"

const {button, div, h3} = van.tags

import { IconError } from "./icons"

const AlertError = ({title="Error", substance}) => {
    const alert = div(
        { class: "p-4 mb-4 text-red-800 border border-red-300 rounded-lg bg-red-50 dark:bg-gray-800 dark:text-red-400 dark:border-red-800" },
    )
    const dismissBtn = button(
        {
            type:"button",
            class: "text-red-800 bg-transparent border border-red-800 hover:bg-red-900 hover:text-white focus:ring-4 focus:outline-none focus:ring-red-300 font-medium rounded-lg text-xs px-3 py-1.5 text-center dark:hover:bg-red-600 dark:border-red-600 dark:text-red-500 dark:hover:text-white dark:focus:ring-red-800",
            onclick: () => alert.remove(),
        },
        "Dismiss"
    )

    const autoDismiss = () => alert.remove();

    setTimeout(autoDismiss, 20000);

    van.add(alert,
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

    return alert
}

export {AlertError}