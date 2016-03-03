import { createAction } from 'redux-actions';
import { WILL_OPEN_SOCKET, DID_OPEN_SOCKET } from '../constants/ActionTypes';

const willOpenSocket = createAction<void>(
  WILL_OPEN_SOCKET
)

const didOpenSocket = createAction<void>(
  DID_OPEN_SOCKET
)

export const openSocket = (socket) => {
  return (dispatch, getState) => {
    dispatch(willOpenSocket())
    return socket.open()
      .then(() => didOpenSocket())
  }
}
