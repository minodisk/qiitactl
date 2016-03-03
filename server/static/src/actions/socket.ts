import { createAction } from 'redux-actions';

import * as types from '../constants/ActionTypes';
import { File } from '../models/file'

const willOpenSocket = createAction<void>(
  types.WILL_OPEN_SOCKET
)

const didOpenSocket = createAction<void>(
  types.DID_OPEN_SOCKET
)

export const openSocket = (socket) => {
  return (dispatch, getState) => {
    dispatch(willOpenSocket())
    return socket.open()
      .then(() => didOpenSocket())
  }
}

const didChangeFile = createAction<File>(
  types.DID_CHANGE_FILE,
  (file:File) => file
)

export const changeFile = (file:File) => {
  return (dispatch, getState) => {
    return dispatch(didChangeFile(file))
  }
}
