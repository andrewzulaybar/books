import { IWork } from './work.model';

export interface IPublication {
  id: number;
  editionPubDate: string;
  format: string;
  imageUrl: string;
  isbn: string;
  isbn13: string;
  language: string;
  numPages: number;
  publisher: string;
  work: IWork;
}
