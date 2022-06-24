package commercetools

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/labd/commercetools-go-sdk/platform"
	"github.com/stretchr/testify/assert"
)

func TestAccStore_createAndUpdateWithID(t *testing.T) {

	name := "test method"
	key := "test-method"
	languages := []string{"en-US"}

	newName := "new test method"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStoreConfig(name, key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"commercetools_store.standard", "name.en", name,
					),
					resource.TestCheckResourceAttr(
						"commercetools_store.standard", "key", key,
					),
					func(s *terraform.State) error {
						res, err := testGetStore(s, "commercetools_store.standard")
						if err != nil {
							return err
						}

						assert.NotNil(t, res)
						assert.EqualValues(t, res.Key, key)
						return nil
					},
				),
			},
			{
				Config: testAccStoreConfig(newName, key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"commercetools_store.standard", "name.en", newName,
					),
					resource.TestCheckResourceAttr(
						"commercetools_store.standard", "key", key,
					),
				),
			},
			{
				Config: testAccStoreConfigWithLanguages(name, key, languages),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"commercetools_store.standard", "languages.#", "1",
					),
					resource.TestCheckResourceAttr(
						"commercetools_store.standard", "languages.0", "en-US",
					),
				),
			},
			{
				Config: testAccNewStoreConfigWithLanguages(name, key, languages),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"commercetools_store.standard", "languages.#", "1",
					),
					resource.TestCheckResourceAttr(
						"commercetools_store.standard", "languages.0", "en-US",
					),
				),
			},
		},
	})
}

func TestAccStore_createAndUpdateDistributionLanguages(t *testing.T) {
	name := "test dl"
	key := "test-dl"
	languages := []string{"en-US"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewStoreConfigWithChannels(name, key, languages),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"commercetools_store.test", "distribution_channels.#", "1",
					),
					resource.TestCheckResourceAttr(
						"commercetools_store.test", "distribution_channels.0", "TEST",
					),
				),
			},
			{
				Config: testAccNewStoreConfigWithoutChannels(name, key, languages),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"commercetools_store.test", "distribution_channels.#", "0",
					),
				),
			},
			{
				Config: testAccNewStoreConfigWithChannels(name, key, languages),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"commercetools_store.test", "distribution_channels.#", "1",
					),
					resource.TestCheckResourceAttr(
						"commercetools_store.test", "distribution_channels.0", "TEST",
					),
				),
			},
		},
	})
}

func TestAccStore_CustomField(t *testing.T) {

	name := "test method"
	key := "standard"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewStoreConfigWithCustomField(name, key, []string{}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"commercetools_store.test", "name.en", name,
					),
					resource.TestCheckResourceAttr(
						"commercetools_store.test", "key", key,
					),
					func(s *terraform.State) error {
						res, err := testGetStore(s, "commercetools_store.test")
						if err != nil {
							return err
						}

						assert.NotNil(t, res)
						assert.NotNil(t, res.Custom)
						assert.NotNil(t, res.Custom.Fields)
						assert.EqualValues(t, res.Custom.Fields["my-field"], "foobar")
						return nil
					},
				),
			},
			{
				Config: testAccNewStoreConfigWithChannels(name, key, []string{}),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						res, err := testGetStore(s, "commercetools_store.test")
						if err != nil {
							return err
						}

						assert.NotNil(t, res)
						assert.Nil(t, res.Custom)
						return nil
					},
				),
			},
		},
	})
}

func testAccStoreConfig(name string, key string) string {
	return fmt.Sprintf(`
	resource "commercetools_store" "standard" {
		name = {
			en = "%[1]s"
			nl = "%[1]s"
		}
		key = "%[2]s"
	}`, name, key)
}

func testAccStoreConfigWithLanguages(name string, key string, languages []string) string {
	return fmt.Sprintf(`
	resource "commercetools_store" "standard" {
		name = {
			en = "%[1]s"
			nl = "%[1]s"
		}
		key = "%[2]s"
		languages = %[3]q
	}`, name, key, languages)
}

func testAccNewStoreConfigWithLanguages(name string, key string, languages []string) string {
	return fmt.Sprintf(`
	resource "commercetools_store" "standard" {
		name = {
			en = "%[1]s"
			nl = "%[1]s"
		}
		key = "%[2]s"
		languages = %[3]q
	}`, name, key, languages)
}

func testAccNewStoreConfigWithChannels(name string, key string, languages []string) string {
	return fmt.Sprintf(`
	resource "commercetools_channel" "test_channel" {
		key = "TEST"
		roles = ["ProductDistribution"]
	}

	resource "commercetools_store" "test" {
		name = {
			en = "%[1]s"
			nl = "%[1]s"
		}
		key = "%[2]s"
		languages = %[3]q
		distribution_channels = [commercetools_channel.test_channel.key]
	}
	`, name, key, languages)
}

func testAccNewStoreConfigWithoutChannels(name string, key string, languages []string) string {
	return fmt.Sprintf(`
	resource "commercetools_channel" "test_channel" {
		key = "TEST"
		roles = ["ProductDistribution"]
	}

	resource "commercetools_store" "test" {
		name = {
			en = "%[1]s"
			nl = "%[1]s"
		}
		key = "%[2]s"
		languages = %[3]q
	}
	`, name, key, languages)
}

func testAccNewStoreConfigWithCustomField(name string, key string, languages []string) string {
	return fmt.Sprintf(`

	resource "commercetools_type" "test" {
		key = "test-for-store"
		name = {
			en = "for Store"
		}
		description = {
			en = "Custom Field for store resource"
		}

		resource_type_ids = ["store"]

		field {
			name = "my-field"
			label = {
				en = "My Custom field"
			}
			type {
				name = "String"
			}
		}
	}

	resource "commercetools_channel" "test_channel" {
		key = "TEST"
		roles = ["ProductDistribution"]
	}

	resource "commercetools_store" "test" {
		name = {
			en = "%[1]s"
			nl = "%[1]s"
		}
		key = "%[2]s"
		languages = %[3]q
		distribution_channels = [commercetools_channel.test_channel.key]
		custom {
			type_id = commercetools_type.test.id
			fields = {
				"my-field" = "foobar"
			}
		}
	}
	`, name, key, languages)
}

func testAccCheckStoreDestroy(s *terraform.State) error {
	client := getClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "commercetools_store":
			{
				response, err := client.Stores().WithId(rs.Primary.ID).Get().Execute(context.Background())
				if err == nil {
					if response != nil && response.ID == rs.Primary.ID {
						return fmt.Errorf("store (%s) still exists", rs.Primary.ID)
					}
					continue
				}

				if newErr := checkApiResult(err); newErr != nil {
					return newErr
				}
			}
		case "commercetools_channel":
			{
				response, err := client.Channels().WithId(rs.Primary.ID).Get().Execute(context.Background())
				if err == nil {
					if response != nil && response.ID == rs.Primary.ID {
						return fmt.Errorf("supply channel (%s) still exists", rs.Primary.ID)
					}
					continue
				}
				if newErr := checkApiResult(err); newErr != nil {
					return newErr
				}
			}
		default:
			continue
		}
	}
	return nil
}

func testGetStore(s *terraform.State, identifier string) (*platform.Store, error) {
	rs, ok := s.RootModule().Resources[identifier]
	if !ok {
		return nil, fmt.Errorf("Store not found")
	}

	client := getClient(testAccProvider.Meta())
	result, err := client.Stores().WithId(rs.Primary.ID).Get().Execute(context.Background())
	if err != nil {
		return nil, err
	}
	return result, nil
}
