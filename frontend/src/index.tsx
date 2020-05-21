import React from 'react';
import ReactDOM from 'react-dom';
import * as serviceWorker from './serviceWorker';
import 'normalize.css';
import 'tailwindcss/dist/tailwind.css';
import { Route, Switch, BrowserRouter as Router } from 'react-router-dom';
import LandingPage from './pages/LandingPage';
import './i18n';

const routing = (
    <Router>
        <Switch>
            <Route exact path="/" component={LandingPage} />
        </Switch>
    </Router>
);

ReactDOM.render(routing, document.getElementById('root'));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
