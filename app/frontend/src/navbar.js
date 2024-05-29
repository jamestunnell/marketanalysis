import van from "vanjs-core"
import { routeTo } from 'vanjs-router'

const {  button, div, li, nav, ul } = van.tags

const NavBarItem = ({text, route, routeArgs, currentRoute}) => {
    const textCls = route === currentRoute ? "text-indigo-500": "hover:text-indigo-400";
    
    return li(
        {class: `md:px-4 md:py-2 ${textCls}`},
        button({onclick: () => { routeTo(route, routeArgs) }}, text),
    )
}

const NavBar = ({currentRoute}) => {
    return div(
        nav(
            {class: "nav font-semibold text-lg bg-gray-200 shadow shadow-gray-300"},
            div(
                {class: "container flex text-gray-500 order-3 w-full"},
                ul(
                    {class: "flex items-center font-semibold flex-wrap"},
                    NavBarItem({text: 'Home', route: 'home', routeArgs: [], currentRoute: currentRoute}),
                    NavBarItem({text: 'Graphs', route: 'graphs', routeArgs: [], currentRoute: currentRoute}),
                )
            )
        )
    )
}

export default NavBar;
