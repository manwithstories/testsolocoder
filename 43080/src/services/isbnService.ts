import type { ISBNBookResult } from '../types/book';
import { logger } from '../utils/logger';
import { validateISBN } from '../utils/validation';

const OPEN_LIBRARY_URL = 'https://openlibrary.org/api/books';
const GOOGLE_BOOKS_URL = 'https://www.googleapis.com/books/v1/volumes';

interface OpenLibraryResponse {
  [key: string]: {
    title?: string;
    authors?: { name: string }[];
    number_of_pages?: number;
    cover?: { small?: string; medium?: string; large?: string };
    subjects?: { name: string }[];
  };
}

interface GoogleBooksResponse {
  items?: Array<{
    volumeInfo: {
      title: string;
      authors?: string[];
      pageCount?: number;
      imageLinks?: { smallThumbnail?: string; thumbnail?: string };
      categories?: string[];
    };
  }>;
}

const fetchFromOpenLibrary = async (isbn: string): Promise<ISBNBookResult | null> => {
  try {
    const url = `${OPEN_LIBRARY_URL}?bibkeys=ISBN:${isbn}&format=json&jscmd=data`;
    const response = await fetch(url);

    if (!response.ok) {
      throw new Error(`OpenLibrary API error: ${response.status}`);
    }

    const data = (await response.json()) as OpenLibraryResponse;
    const key = `ISBN:${isbn}`;
    const bookData = data[key];

    if (!bookData) {
      logger.info(`No book found in OpenLibrary for ISBN: ${isbn}`);
      return null;
    }

    return {
      title: bookData.title || '',
      author: bookData.authors?.map(a => a.name).join(', ') || '',
      isbn,
      coverUrl: bookData.cover?.large || bookData.cover?.medium || bookData.cover?.small,
      totalPages: bookData.number_of_pages,
      categories: bookData.subjects?.map(s => s.name).slice(0, 5),
    };
  } catch (error) {
    logger.error('Failed to fetch from OpenLibrary', error);
    return null;
  }
};

const fetchFromGoogleBooks = async (isbn: string): Promise<ISBNBookResult | null> => {
  try {
    const url = `${GOOGLE_BOOKS_URL}?q=isbn:${isbn}`;
    const response = await fetch(url);

    if (!response.ok) {
      throw new Error(`Google Books API error: ${response.status}`);
    }

    const data = (await response.json()) as GoogleBooksResponse;
    const bookData = data.items?.[0]?.volumeInfo;

    if (!bookData) {
      logger.info(`No book found in Google Books for ISBN: ${isbn}`);
      return null;
    }

    return {
      title: bookData.title,
      author: bookData.authors?.join(', ') || '',
      isbn,
      coverUrl: bookData.imageLinks?.thumbnail?.replace('http://', 'https://'),
      totalPages: bookData.pageCount,
      categories: bookData.categories,
    };
  } catch (error) {
    logger.error('Failed to fetch from Google Books', error);
    return null;
  }
};

export const fetchBookByISBN = async (isbn: string): Promise<ISBNBookResult | null> => {
  const cleanedISBN = isbn.replace(/[-\s]/g, '');

  if (!validateISBN(cleanedISBN)) {
    logger.warn(`Invalid ISBN format: ${isbn}`);
    throw new Error('ISBN格式不正确');
  }

  logger.info(`Fetching book info for ISBN: ${cleanedISBN}`);

  const results = await Promise.allSettled([
    fetchFromOpenLibrary(cleanedISBN),
    fetchFromGoogleBooks(cleanedISBN),
  ]);

  for (const result of results) {
    if (result.status === 'fulfilled' && result.value) {
      logger.info('Book info fetched successfully', result.value);
      return result.value;
    }
  }

  logger.warn(`No book info found for ISBN: ${cleanedISBN}`);
  return null;
};
