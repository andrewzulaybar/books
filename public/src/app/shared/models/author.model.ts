import { ILocation } from './location.model';

export interface IAuthor {
  id: number;
  firstName: string;
  lastName: string;
  gender: string;
  dateOfBirth: string;
  placeOfBirth: ILocation;
}
