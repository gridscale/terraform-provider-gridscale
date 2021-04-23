package gridscale

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/gridscale/gsclient-go/v3"
)

const (
	privateKeyExample = "-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEA0gzQuib+IjCugSwGtLOaJt9VGGkj0FAzXH0CKASGOje21v3T\n9M2luERn4SwSVVVAmBvLLp7SMhWKTGxGZ37wodywYQiJJViTbDRWSHEmqld09MxU\nA0vwbx9KUJ5SwZUKno9WIqFIG30TKm2pRJ7XupWCKBdpzjq7FpyRZ1KM24S2fDea\n4IzuHRwJvez2B5EVo9s69Ng1Wb/AWLkjSOhB4lhdXgsqNPfeDZ+IqFvpzhvh8+6L\nkVhSFImAHijfm4VH+oGyi1rfw5nCXh1QgE+KtvyYVPpu5jAP93DYNVPGMfFoPxpa\n1BaDz3pTWxz8DcI+9Zv4TIP/gmfdjMq851fs9QIDAQABAoIBAQCC0VhtlztPL+J+\nsFy8zd6DtzPYMCMN6eUTBQgCxE6Uk35cB1mqbSl4yaGoqzn/y0ifcy/nqSJf7upj\nDkwKyRQ0uq1ewWIo0DGAglf6Zf4euLkufNMV+oWTr3kpHHGN/rTOa5bJmoi+FfeJ\nxSJWilbw7YxSS8KpngpG2tvz7YWscKmUEYCawj5NO+LjN2CuBWoaoogz098Cj6Ds\nhyZMgRcW15WnFRqLb9zU4E6+IuT6ubwUvOm394CX2iWQxf+M+1+gYkHE9dTQCBcE\nVTyYNzZFZxjPJ3QKUH4FQHXBvaKoHurRyjZ8wePRDXfP9k7RAh0JCNciXTcDkR1G\nrp23WRQBAoGBAPaFMxJYa4COagEQfRNwmnjwS0OacVzx6i9JJ03LCdQPVO96AOdy\n4W0fh5Z9oOK24Y9E9rHkftIgsj7kL31zEKcT6gDUnyR+Z3nNVogL6DN+ee7o6T1Y\nFio/8wjLVYxls9NUnh7aO9VL+jWOThbhxhdeg+8MYlZvfrVAGpFVbJYFAoGBANog\nmDktF7y9MfujIv2M/XxJgrOs6pc8fYkOg/ewZnwf/gg944Ou2uFLo+pZYmWJo1mE\nXLbsKiE5N7XEltW+6pB4MVIh9noaxECo/Co23RBklO/uUs3kxw3SZUFM+EcCHsFE\neF50IbpWO7h1trfENt+RPJBBTvYKIQxz864Xcj4xAoGBALsrsPYp2h12g7ehLEX5\nRd+zGoIOninip1sVoNjo4uJ2bZeRu2P/jnbK7KhvZBsena2ov+GbSksCdjNYD1Kh\nbbBdZ1sQFHCh5LDPHfcXi0qlqbCieIyR9+pGOaGguKhrg5hd89wa+Djfge7tyA6J\nh1yw710tuvBDoSG+tS1JfbN9AoGAZvKpRZr9sAuPUpeNRSPM42MCd/Cti5F2IkCi\nkj1EMWMcDyG+hoD5NwnvJATpxshuOyqYY5EU1TdfDr/o/JqttF+1z3DOFpUZePcI\nhJf+P+Vo0dCmbhyCwghB188PAtyCTW8WbGXdIUffydBs21CVIXE3p1HbhY4ESQ9U\nvB6+wZECgYEArEk/WbFnosDo9Tj4EcJbl5bZg0CDq5dTF3W90+2hLzpdiZiRQknM\n7xH3UQ2AjjLgt944uMR4P1l1bNDG7oHswnPb4O3IrlzQYUpaizRzhnuRYUPkokhr\ng4njwtk/2O7ytXSa1ewR7+yv+U5BiOTVWoZfWlAxDA6VUgQ8Nuzj0CA=\n-----END RSA PRIVATE KEY-----"
	leafCertExample   = "-----BEGIN CERTIFICATE-----\nMIIC1TCCAb2gAwIBAgIJAP4/xUQz5NypMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNV\nBAMTD3d3dy5leGFtcGxlLmNvbTAeFw0yMTAzMjkxNDQxNTdaFw0zMTAzMjcxNDQx\nNTdaMBoxGDAWBgNVBAMTD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEB\nBQADggEPADCCAQoCggEBANIM0Lom/iIwroEsBrSzmibfVRhpI9BQM1x9AigEhjo3\nttb90/TNpbhEZ+EsElVVQJgbyy6e0jIVikxsRmd+8KHcsGEIiSVYk2w0VkhxJqpX\ndPTMVANL8G8fSlCeUsGVCp6PViKhSBt9EyptqUSe17qVgigXac46uxackWdSjNuE\ntnw3muCM7h0cCb3s9geRFaPbOvTYNVm/wFi5I0joQeJYXV4LKjT33g2fiKhb6c4b\n4fPui5FYUhSJgB4o35uFR/qBsota38OZwl4dUIBPirb8mFT6buYwD/dw2DVTxjHx\naD8aWtQWg896U1sc/A3CPvWb+EyD/4Jn3YzKvOdX7PUCAwEAAaMeMBwwGgYDVR0R\nBBMwEYIPd3d3LmV4YW1wbGUuY29tMA0GCSqGSIb3DQEBBQUAA4IBAQDKWvBcVS1R\n8K+HYBpEaVnElMk7vBOvdq7h/SZSXVAGNgNM1RPK7w6zWyDwI2Zs5COQiAbK0MwL\nBkI6RwTfDd8RL/nZe35iL6agI0CQbX3/l7Zo09n2RpShpzHbzWIkuzPNlzR+b0pb\nIGyChAsL1O+d3Ft/8LFkFiXIcEb+0kB75X/R+Tx2+LmCMxCHM4JYudzE8mzuKHuR\neXktzUNvVHzkcAN1rpTCK0tGJivmh+QMrwXZh2cWsL6xdI8m24J+e2zJef/mosYQ\n+xX/qCCQKbsRfckG3SDHi6RZlCnzrwSaW+djekAIuju+3HIbeDjjAI4DUSfqd1SH\n87fnnaj9tE69\n-----END CERTIFICATE-----"
)

func TestAccResourceGridscaleSSLCert_Basic(t *testing.T) {
	var object gsclient.SSLCertificate
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGridscaleSSLCertDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceGridscaleSSLCertConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGridscaleSSLCertExists("gridscale_ssl_certificate.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_ssl_certificate.foo", "name", name),
					resource.TestCheckResourceAttr(
						"gridscale_ssl_certificate.foo", "private_key", privateKeyExample),
					resource.TestCheckResourceAttr(
						"gridscale_ssl_certificate.foo", "leaf_certificate", leafCertExample),
					resource.TestCheckResourceAttrSet("gridscale_ssl_certificate.foo", "id"),
					resource.TestCheckResourceAttrSet("gridscale_ssl_certificate.foo", "common_name"),
					resource.TestCheckResourceAttrSet("gridscale_ssl_certificate.foo", "status"),
					resource.TestCheckResourceAttrSet("gridscale_ssl_certificate.foo", "create_time"),
					resource.TestCheckResourceAttrSet("gridscale_ssl_certificate.foo", "change_time"),
					resource.TestCheckResourceAttrSet("gridscale_ssl_certificate.foo", "not_valid_after"),
					resource.TestCheckResourceAttr("gridscale_ssl_certificate.foo", "fingerprints.#", "1"),
					resource.TestCheckResourceAttrSet("gridscale_ssl_certificate.foo", "fingerprints.0.md5"),
					resource.TestCheckResourceAttrSet("gridscale_ssl_certificate.foo", "fingerprints.0.sha256"),
					resource.TestCheckResourceAttrSet("gridscale_ssl_certificate.foo", "fingerprints.0.sha1"),
				),
			},
		},
	})
}

func testAccCheckResourceGridscaleSSLCertExists(n string, object *gsclient.SSLCertificate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No object UUID is set")
		}

		client := testAccProvider.Meta().(*gsclient.Client)

		id := rs.Primary.ID

		foundObject, err := client.GetSSLCertificate(context.Background(), id)

		if err != nil {
			return err
		}

		if foundObject.Properties.ObjectUUID != id {
			return fmt.Errorf("Object not found")
		}

		*object = foundObject

		return nil
	}
}

func testAccCheckGridscaleSSLCertDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*gsclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gridscale_ssl_certificate" {
			continue
		}

		_, err := client.GetSSLCertificate(context.Background(), rs.Primary.ID)
		if err != nil {
			if requestError, ok := err.(gsclient.RequestError); ok {
				if requestError.StatusCode != 404 {
					return fmt.Errorf("Object %s still exists", rs.Primary.ID)
				}
			} else {
				return fmt.Errorf("Unable to fetching object %s", rs.Primary.ID)
			}
		} else {
			return fmt.Errorf("Object %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckResourceGridscaleSSLCertConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_ssl_certificate" "foo" {
  name   = "%s"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEA0gzQuib+IjCugSwGtLOaJt9VGGkj0FAzXH0CKASGOje21v3T\n9M2luERn4SwSVVVAmBvLLp7SMhWKTGxGZ37wodywYQiJJViTbDRWSHEmqld09MxU\nA0vwbx9KUJ5SwZUKno9WIqFIG30TKm2pRJ7XupWCKBdpzjq7FpyRZ1KM24S2fDea\n4IzuHRwJvez2B5EVo9s69Ng1Wb/AWLkjSOhB4lhdXgsqNPfeDZ+IqFvpzhvh8+6L\nkVhSFImAHijfm4VH+oGyi1rfw5nCXh1QgE+KtvyYVPpu5jAP93DYNVPGMfFoPxpa\n1BaDz3pTWxz8DcI+9Zv4TIP/gmfdjMq851fs9QIDAQABAoIBAQCC0VhtlztPL+J+\nsFy8zd6DtzPYMCMN6eUTBQgCxE6Uk35cB1mqbSl4yaGoqzn/y0ifcy/nqSJf7upj\nDkwKyRQ0uq1ewWIo0DGAglf6Zf4euLkufNMV+oWTr3kpHHGN/rTOa5bJmoi+FfeJ\nxSJWilbw7YxSS8KpngpG2tvz7YWscKmUEYCawj5NO+LjN2CuBWoaoogz098Cj6Ds\nhyZMgRcW15WnFRqLb9zU4E6+IuT6ubwUvOm394CX2iWQxf+M+1+gYkHE9dTQCBcE\nVTyYNzZFZxjPJ3QKUH4FQHXBvaKoHurRyjZ8wePRDXfP9k7RAh0JCNciXTcDkR1G\nrp23WRQBAoGBAPaFMxJYa4COagEQfRNwmnjwS0OacVzx6i9JJ03LCdQPVO96AOdy\n4W0fh5Z9oOK24Y9E9rHkftIgsj7kL31zEKcT6gDUnyR+Z3nNVogL6DN+ee7o6T1Y\nFio/8wjLVYxls9NUnh7aO9VL+jWOThbhxhdeg+8MYlZvfrVAGpFVbJYFAoGBANog\nmDktF7y9MfujIv2M/XxJgrOs6pc8fYkOg/ewZnwf/gg944Ou2uFLo+pZYmWJo1mE\nXLbsKiE5N7XEltW+6pB4MVIh9noaxECo/Co23RBklO/uUs3kxw3SZUFM+EcCHsFE\neF50IbpWO7h1trfENt+RPJBBTvYKIQxz864Xcj4xAoGBALsrsPYp2h12g7ehLEX5\nRd+zGoIOninip1sVoNjo4uJ2bZeRu2P/jnbK7KhvZBsena2ov+GbSksCdjNYD1Kh\nbbBdZ1sQFHCh5LDPHfcXi0qlqbCieIyR9+pGOaGguKhrg5hd89wa+Djfge7tyA6J\nh1yw710tuvBDoSG+tS1JfbN9AoGAZvKpRZr9sAuPUpeNRSPM42MCd/Cti5F2IkCi\nkj1EMWMcDyG+hoD5NwnvJATpxshuOyqYY5EU1TdfDr/o/JqttF+1z3DOFpUZePcI\nhJf+P+Vo0dCmbhyCwghB188PAtyCTW8WbGXdIUffydBs21CVIXE3p1HbhY4ESQ9U\nvB6+wZECgYEArEk/WbFnosDo9Tj4EcJbl5bZg0CDq5dTF3W90+2hLzpdiZiRQknM\n7xH3UQ2AjjLgt944uMR4P1l1bNDG7oHswnPb4O3IrlzQYUpaizRzhnuRYUPkokhr\ng4njwtk/2O7ytXSa1ewR7+yv+U5BiOTVWoZfWlAxDA6VUgQ8Nuzj0CA=\n-----END RSA PRIVATE KEY-----"
  leaf_certificate = "-----BEGIN CERTIFICATE-----\nMIIC1TCCAb2gAwIBAgIJAP4/xUQz5NypMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNV\nBAMTD3d3dy5leGFtcGxlLmNvbTAeFw0yMTAzMjkxNDQxNTdaFw0zMTAzMjcxNDQx\nNTdaMBoxGDAWBgNVBAMTD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEB\nBQADggEPADCCAQoCggEBANIM0Lom/iIwroEsBrSzmibfVRhpI9BQM1x9AigEhjo3\nttb90/TNpbhEZ+EsElVVQJgbyy6e0jIVikxsRmd+8KHcsGEIiSVYk2w0VkhxJqpX\ndPTMVANL8G8fSlCeUsGVCp6PViKhSBt9EyptqUSe17qVgigXac46uxackWdSjNuE\ntnw3muCM7h0cCb3s9geRFaPbOvTYNVm/wFi5I0joQeJYXV4LKjT33g2fiKhb6c4b\n4fPui5FYUhSJgB4o35uFR/qBsota38OZwl4dUIBPirb8mFT6buYwD/dw2DVTxjHx\naD8aWtQWg896U1sc/A3CPvWb+EyD/4Jn3YzKvOdX7PUCAwEAAaMeMBwwGgYDVR0R\nBBMwEYIPd3d3LmV4YW1wbGUuY29tMA0GCSqGSIb3DQEBBQUAA4IBAQDKWvBcVS1R\n8K+HYBpEaVnElMk7vBOvdq7h/SZSXVAGNgNM1RPK7w6zWyDwI2Zs5COQiAbK0MwL\nBkI6RwTfDd8RL/nZe35iL6agI0CQbX3/l7Zo09n2RpShpzHbzWIkuzPNlzR+b0pb\nIGyChAsL1O+d3Ft/8LFkFiXIcEb+0kB75X/R+Tx2+LmCMxCHM4JYudzE8mzuKHuR\neXktzUNvVHzkcAN1rpTCK0tGJivmh+QMrwXZh2cWsL6xdI8m24J+e2zJef/mosYQ\n+xX/qCCQKbsRfckG3SDHi6RZlCnzrwSaW+djekAIuju+3HIbeDjjAI4DUSfqd1SH\n87fnnaj9tE69\n-----END CERTIFICATE-----"
}
`, name)
}
