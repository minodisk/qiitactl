/// <reference path="../../typings/tsd.d.ts" />

import * as React from 'react'
import Link from './Link'

const PropTypes = React.PropTypes

const File = ({onClick, path}) => (
  <li
    onClick={onClick}
  >
    <Link
      active={true}
      onClick={() => onClick()}
      children={path}
    />
  </li>
);

File['propTypes'] = {
  onClick: PropTypes.func.isRequired,
  path: PropTypes.string.isRequired,
}

export default File
