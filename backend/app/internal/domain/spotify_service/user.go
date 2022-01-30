package spotify_service

import "fmt"

func (c *Client) GetCurrentUserId() (string, error) {
	user, err := c.client.CurrentUser()
	if err != nil {
		err = fmt.Errorf("unable to get current user: %w", err)
		return "", err
	}
	return user.ID, nil
}
