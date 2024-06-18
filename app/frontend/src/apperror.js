import van from "vanjs-core"

import {AlertError} from './alerts.js'
import { ModalBackground, ModalForeground } from './modal.js'

const { div, li, p, ul } = van.tags

const AppErrorContent = (appErr) => {
    return div(
        {class: "flex flex-col space-y-2"},
        p(`Message: ${appErr.message}`),
        () => appErr.details ? div(
            p(`Details:`),
            ul(
                {class: "list-disc list-inside"},
                appErr.details.map(detail => li(detail))
            ),
        ) : null,
    )
}

const AppErrorAlert = (appErr) => {
    console.log("adding app error alert", appErr)

    const alert = new AlertError({title: appErr.title, substance: AppErrorContent(appErr)})

    van.add(document.body, alert.render());
}

const AppErrorModal = (appErr) => {
    console.log("showing app error modal", appErr)

    const modal = ModalBackground(ModalForeground({}, AppErrorContent(appErr)))

    van.add(document.body, modal);
}

export {AppErrorAlert, AppErrorModal};