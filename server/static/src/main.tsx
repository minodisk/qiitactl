import * as React from 'react';
import {render} from 'react-dom';
import {createHistory} from 'history'
import {Route} from 'react-router'
import {
  connect,
  Provider
} from 'react-redux';
import {
  Store,
  compose,
  createStore,
  bindActionCreators,
  combineReducers,
  applyMiddleware
} from 'redux';
import {
  ReduxRouter,
  reduxReactRouter
} from 'redux-router'
import * as thunkMiddleware from 'redux-thunk'
const createLogger = require('redux-logger')

import { rootReducer } from './reducers/rootReducer';
import { openSocket  } from './actions/socket'
import { fetchFiles } from './actions/files'
import App from './containers/App';

const routes = (
  <Route path="/" component={App}>
  </Route>
)

const loggerMiddleware = createLogger()
const store: Store = compose(
  applyMiddleware(
    thunkMiddleware,
    loggerMiddleware
  ),
  reduxReactRouter({
    routes,
    createHistory
  })
)(createStore)(rootReducer);

store
  .dispatch(openSocket())
  .then(() => store.dispatch(fetchFiles()))
  .then(() => console.log('done:', store.getState()))

render(
  <Provider store={store}>
    <ReduxRouter/>
  </Provider>,
  document.querySelector('#app')
);
