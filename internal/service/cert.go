package service

import "net/http"

type CertService struct {
}

func NewCertService() *CertService {
	return &CertService{}
}

func (s *CertService) CAProviders(w http.ResponseWriter, r *http.Request) {

}

func (s *CertService) DNSProviders(w http.ResponseWriter, r *http.Request) {

}

func (s *CertService) Algorithms(w http.ResponseWriter, r *http.Request) {

}

func (s *CertService) List(w http.ResponseWriter, r *http.Request) {

}

func (s *CertService) Create(w http.ResponseWriter, r *http.Request) {

}

func (s *CertService) Update(w http.ResponseWriter, r *http.Request) {

}

func (s *CertService) Get(w http.ResponseWriter, r *http.Request) {

}

func (s *CertService) Delete(w http.ResponseWriter, r *http.Request) {

}

func (s *CertService) Obtain(w http.ResponseWriter, r *http.Request) {

}

func (s *CertService) Renew(w http.ResponseWriter, r *http.Request) {

}

func (s *CertService) ManualDNS(w http.ResponseWriter, r *http.Request) {

}

func (s *CertService) Deploy(w http.ResponseWriter, r *http.Request) {

}
