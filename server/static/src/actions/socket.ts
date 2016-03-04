import { createAction } from 'redux-actions';

import * as types from '../constants/ActionTypes';
import { Connection } from '../models/socket'

const willOpenSocket = createAction<Connection>(
  types.WILL_OPEN_SOCKET,
  () => ({opened: false})
)

const didOpenSocket = createAction<Connection>(
  types.DID_OPEN_SOCKET,
  () => ({opened: true})
)

export const openSocket = (socket) => {
  return (dispatch, getState) => {
    dispatch(willOpenSocket())
    return socket.open()
      .then(() => dispatch(didOpenSocket()))
  }
}

export const didCloseSocket = createAction<Connection>(
  types.DID_CLOSE_SOCKET,
  () => ({opened: false})
)

const didChangeFile = createAction<File>(
  types.DID_CHANGE_FILE,
  (file:File) => file
)

export const changeFile = (file:File) => {
  return (dispatch, getState) => {
    return dispatch(didChangeFile(file))
  }
}
