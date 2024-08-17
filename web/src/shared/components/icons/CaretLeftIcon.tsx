import React, { ReactElement } from 'react';
import caretLeft from '../../../assets/icons/caret-left.svg';

const CaretLeftIcon = (props: React.ImgHTMLAttributes<HTMLImageElement>): ReactElement => <img src={caretLeft} alt='caret-left' {...props} />

export default CaretLeftIcon