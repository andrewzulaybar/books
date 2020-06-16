import { IAuthor } from './author.model';

export interface IWork {
  id: number;
  description: string;
  initialPubDate: string;
  originalLanguage: string;
  title: string;
  author: IAuthor;
}
