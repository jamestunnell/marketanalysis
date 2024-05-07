import van from "vanjs-core"
import {Modal} from "vanjs-ui"

import AppError from "./apperror.js"
import {ButtonAct, ButtonCancel} from './buttons.js'

const {div, h3, input, option, p, select, table, tbody, td, th, thead, tr} = van.tags

const BASE_URL = `http://localhost:4002`
const TIME_ZONES = Intl.supportedValuesOf('timeZone');
const TZ_NEW_YORK = "America/New_York";

const getSecurities = async () => {
    console.log("getting securities");

    const resp = await fetch(`${BASE_URL}/securities`, {credentials: 'same-origin'});

    if (!resp.ok) {
        console.log("non-ok get securities response %d: %s", resp.status, resp.text());

        return []
    }
    
    const d = await resp.json()

    console.log("get securities response JSON: %o", d)

    return d.securities;
}

const addSecurity = async (item) => {
    return await fetch(`${BASE_URL}/securities`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json; charset=utf-8'
        },
        body: JSON.stringify(item),
        credentials: 'same-origin'
    });
}

const delSecurity = async (symbol) => {
    console.log("deleting security");

    const resp = await fetch(`${BASE_URL}/securities/${symbol}`, {
        method: 'DELETE',
        credentials: 'same-origin'
    });

    console.log('delete security result:', resp.status);

    const success = resp.status === 200;
    
    if (!success) {
        console.log(resp);
    }

    return success; 
}

const AddSecurityInModal = ({tableBody, modalClosed}) => {
    const sym = van.state("");
    const tz = van.state(TZ_NEW_YORK);
    const open = van.state("09:30");
    const close = van.state("16:00");
    const errType = van.state("");
    const errMsg = van.state("");
    const errDetails = van.state([]);
    const errIsVisible = van.state(false)

    return div(
        {class: "space-y-6"},
        h3({class: "text-2xl font-bold"}, "Add New Security"),
        div(
            {class: "grid grid-cols-2 gap-4"},
            p({class: "text-right"}, "Symbol"),
            input({type: "text", value: sym, oninput: e => sym.val = e.target.value, placeholder: "SPY, QQQ, etc."}),
            p({class: "text-right"}, "Time Zone"),
            select(
                {oninput: (e) => tz.val = e.target.value},
                TIME_ZONES.map(x => {
                    let props = {value: x};
                    if (x === TZ_NEW_YORK) {
                        props.selected = "selected";
                    }
                    
                    return option(props, x)                    
                }),
            ),
            p({class: "text-right"}, "Open"),
            input({type: "text", value: open, oninput: e => open.val = e.target.value}),
            p({class: "text-right"}, "Close"),
            input({type: "text", value: close, oninput: e => close.val = e.target.value}),
        ),
        AppError({type: errType, msg: errMsg, details: errDetails, isVisible: errIsVisible}),
        div(
            {class:"mt-4 flex justify-end"},
            ButtonCancel({text: "Cancel", onclick: () => modalClosed.val = true}),
            ButtonAct({text: "OK", onclick: async () => {
                const item = {symbol: sym.val, timeZone: tz.val, open: open.val, close: close.val};
                const resp = await addSecurity(item);

                if (resp.status !== 204) {
                    console.log('failed to add security', resp);

                    const err = await resp.json();
                    
                    errType.val = err.errType;
                    errMsg.val = err.message;
                    errDetails.val = err.details;
                    errIsVisible.val = true;
                    
                    return
                }
            
                console.log('added security %s', item);
            
                errIsVisible.val = false;
                modalClosed.val = true;

                van.add(tableBody, TableRow(item))
            }}),
        ),

    )
}

const TableRow = (item) => {
    const deleted = van.state(false)
    return () => deleted.val ? null : tr(
        {class: "border border-solid"},
        td({class: "px-6 py-4 font-medium"}, item.symbol),
        td({class: "px-6 py-4"}, item.timeZone),
        td({class: "px-6 py-4"}, item.open),
        td({class: "px-6 py-4"}, item.close),
        td({class: "px-6 py-4"}, ButtonAct({
                text: "Delete",
                onclick: () => {
                    if (delSecurity(item.symbol)) {
                        deleted.val = true
                    }
                },
            }),
        ),
    )
}

const Securities = () => {
    const tableBody = tbody({class:"table-auto"})

    getSecurities().then(
        (items) => {
            console.log("found %d securities", items.length);
            
            const rows = items.map(item => TableRow(item));
            
            van.add(tableBody, rows)
        }
    ).catch(error => {
        console.log("failed to get securities: " + error);
    });

    return div(
        {class: "w-full p-4 space-y-6"},
        table(
            {class:"table-auto min-w-full text-left text-sm"},
            thead(
                tr(
                    {class: "border-b border-neutral-200 font-medium"},
                    th({class: "px-6 py-4"}, "Symbol"),
                    th({class: "px-6 py-4"}, "Time Zone"),
                    th({class: "px-6 py-4"}, "Open"),
                    th({class: "px-6 py-4"}, "Close"),
                    th({class: "px-6 py-4"}, ""),
                ),
            ),
            tableBody
        ),
        ButtonAct({text: "Add New", onclick: () => {
            const closed = van.state(false)
            
            van.add(document.body, Modal({closed},
              div({style: "display: flex; justify-content: center;"},
                AddSecurityInModal({tableBody: tableBody, modalClosed: closed}),
              ),
            ))
        }}),
    );
}

export default Securities