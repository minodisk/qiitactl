import * as React from 'react';
import { render } from 'react-dom';
import { createHistory } from 'history'
import {
  Route,
  IndexRoute,
} from 'react-router'
import {
  Provider
} from 'react-redux';
import {
  Store,
  compose,
  applyMiddleware,
  createStore,
} from 'redux';
import {
  ReduxRouter,
  reduxReactRouter
} from 'redux-router'
import * as thunkMiddleware from 'redux-thunk'
const createLogger = require('redux-logger')

import { rootReducer } from './reducers/rootReducer';
import { Socket } from './net/socket'
import { openSocket } from './actions/socket'
import { getTree } from './actions/tree'
import App from './containers/App';
import File from './containers/File'
import NotFound from './containers/NotFound'

const routes = (
  <Route path="/" component={App}>
    <IndexRoute component={File}/>
    <Route path="/**/*.md" component={File}/>
    <Route path="*" component={NotFound}/>
  </Route>
)

const loggerMiddleware = createLogger()
const store:Store = compose(
  applyMiddleware(
    thunkMiddleware,
    loggerMiddleware
  ),
  reduxReactRouter({
    routes,
    createHistory
  })
)(createStore)(rootReducer);

const socket: Socket = Socket.getInstance()
socket.dispatch = store.dispatch

store
  .dispatch(openSocket(socket))
  .then(() => {
    store.dispatch(getTree())
    render(
      <Provider store={store}>
        <ReduxRouter/>
      </Provider>,
      document.querySelector('#app')
    );
  })
  .catch(err => console.error(err))
