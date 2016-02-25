import * as React from 'react';

interface LinkProps {
  indent: number;
  className: string;
  active: boolean
  children: any
  onClick: Function
}

export default class Link extends React.Component<LinkProps, void> {
  render() {
    const {indent, className, active, children, onClick} = this.props
    if (!active) {
      return (
        <span>{children}</span>
      )
    }
    return (
      <a
        className={className}
        style={{
          paddingLeft: (1.5 * indent) + 'em'
        }}
        href="#"
        onClick={(e) => {
          e.preventDefault()
          onClick(e)
        }}>
        {children}
      </a>
    )
  }
}
