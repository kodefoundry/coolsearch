import React from 'react';
import logo from './logo.svg';
import './css/App.css';
import Homepage from "./components/HomePage";
import Searchresult from './components/SearchResult';
import {library} from '@fortawesome/fontawesome-svg-core'
import {faTrash} from '@fortawesome/free-solid-svg-icons'
import {BrowserRouter as Router, Route, Link} from 'react-router-dom';
import 'bootstrap/dist/css/bootstrap.min.css';

library.add(faTrash)

class App extends React.Component {
    render() {
        return (
            <Router>
                <div className="AppBg">
                  <Route exact path='/' component={Homepage}/>
                  <Route path='/Result/:uuid' component={Searchresult}/>
                </div>
            </Router>
        );
    }
}

export default App;
