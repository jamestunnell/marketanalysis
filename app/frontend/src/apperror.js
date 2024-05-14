import van from "vanjs-core"
import {Modal} from "vanjs-ui"

import { Button } from "./buttons"

const { div, p } = van.tags

const DoAppErrorModal = (appErr) => {
    const closed = van.state(false)
    
    van.add(
        document.body,
        Modal({closed},
            div(
                {class: "flex flex-col drop-shadow hover:drop-shadow-lg w-300 rounded-md"},
                p({class: "text-2xl font-medium font-bold text-center"}, `Error: ${appErr.title}`),
                p({class: "text-lg font-medium"}, appErr.message),
                div(
                    {class: "flex flex-col overflow-auto p-3"},
                    appErr.details.map(detail => p(detail)),
                ),
                div(
                    {class:"mt-4 flex justify-center"},
                    Button({child: "OK", onclick: ()=>closed.val = true}),
                ),
            )
        ),
    );
}

export {DoAppErrorModal};