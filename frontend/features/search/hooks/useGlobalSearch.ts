"use client";

import {
  useEffect,
  useState,
} from "react";

import { searchCRM } from "../api";
import { SearchResponse } from "../types";

export default function useGlobalSearch(
  query: string
) {
  const [loading, setLoading] =
    useState(false);

  const [results, setResults] =
    useState<SearchResponse>({
      leads: [],
      contacts: [],
      tasks: [],
    });

  useEffect(() => {
    if (query.trim().length < 2) {
      setResults({
        leads: [],
        contacts: [],
        tasks: [],
      });

      return;
    }

    const timeout = setTimeout(async () => {
      try {
        setLoading(true);

        const data =
          await searchCRM(query);

        setResults(data);
      } catch {
        setResults({
          leads: [],
          contacts: [],
          tasks: [],
        });
      } finally {
        setLoading(false);
      }
    }, 350);

    return () =>
      clearTimeout(timeout);
  }, [query]);

  return {
    loading,
    results,
  };
}