import { handleActions, Action } from 'redux-actions';

import { Connection } from '../models/socket'
import {
  DID_OPEN_SOCKET,
  DID_CLOSE_SOCKET,
} from '../constants/ActionTypes';

export default handleActions<Connection>({
  [DID_OPEN_SOCKET]: (state:Connection, action:Action):Connection => {
    return action.payload
  },
  [DID_CLOSE_SOCKET]: (state:Connection, action:Action):Connection => {
    return action.payload
  },
}, null);
