import * as React from 'react';
import { Link } from 'react-router'

interface Props {
  indent: number;
  className: string;
  active: boolean
  title: string
  children: any
  linkTo?: string
  onClick?: Function
}

export default class TreeElement extends React.Component<Props, void> {
  render() {
    const {indent, className, active, title, children, linkTo, onClick} = this.props
    if (!active) {
      return (
        <span>{children}</span>
      )
    }

    if (linkTo != null) {
      return (
        <Link
          className={className}
          style={{
            paddingLeft: (1.5 * indent) + 'em'
          }}
          title={title}
          to={'/markdown/' + linkTo}
        >
          {children}
        </Link>
      )
    }

    return (
      <a
        className={className}
        style={{
          paddingLeft: (1.5 * indent) + 'em'
        }}
        href="#"
        title={title}
        onClick={(e) => {
          e.preventDefault()
          onClick(e)
        }}
      >
        {children}
      </a>
    )
  }
}
