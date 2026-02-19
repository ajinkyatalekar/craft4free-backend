package service

import (
    "github.com/supabase-community/supabase-go"
)

func NewSupabaseClient(url, key string) (*supabase.Client, error) {
    client, err := supabase.NewClient(url, key, &supabase.ClientOptions{})
    if err != nil {
        return nil, err
    }
    
    return client, nil
}