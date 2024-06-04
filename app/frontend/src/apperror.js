import van from "vanjs-core"

import {AlertError} from './alerts.js'
import capitalize from "./capitalize.js"

const { div, li, p, ul } = van.tags

// const DoAppErrorModal = (appErr) => {
//     const closed = van.state(false)
    
//     van.add(
//         document.body,
//         Modal({closed},
//             div(
//                 {class: "flex flex-col drop-shadow hover:drop-shadow-lg w-300 rounded-md"},
//                 p({class: "text-2xl font-medium font-bold text-center"}, `Error: ${appErr.title}`),
//                 p({class: "text-lg font-medium"}, appErr.message),
//                 div(
//                     {class: "flex flex-col overflow-auto p-3"},
//                     appErr.details.map(detail => p(detail)),
//                 ),
//                 div(
//                     {class:"mt-4 flex justify-center"},
//                     Button({child: "OK", onclick: ()=>closed.val = true}),
//                 ),
//             )
//         ),
//     );
// }

const AppErrorAlert = (appErr) => {
    console.log("adding app error alert", appErr)

    const substance = div(
        {class: "flex flex-col space-y-2"},
        p(`Message: ${appErr.message}`),
        div(
            p(`Details:`),
            ul(
                {class: "list-disc list-inside"},
                appErr.details.map(detail => li(detail))
            ),
        )
    )

    const alert = new AlertError({title: appErr.title, substance: substance})

    van.add(document.body, alert.render());
}

export {AppErrorAlert};