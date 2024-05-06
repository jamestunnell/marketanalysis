import van from "vanjs-core"
import { Route } from 'vanjs-router'
import Securities from './src/securities.js'
import Graphs from './src/graphs.js'
import NavBar from './src/navbar.js'

const { div } = van.tags

const App = () => {
    return div(
        Route({ name: 'home' },
            NavBar({currentRoute: 'home'}),
        ),
        Route({ name: 'securities' },
            NavBar({currentRoute: 'securities'}),
            Securities(),
        ),
        Route({ name: 'graphs' },
            NavBar({currentRoute: 'graphs'}),
            Graphs(),
        ),
    )
}

van.add(document.body, App())