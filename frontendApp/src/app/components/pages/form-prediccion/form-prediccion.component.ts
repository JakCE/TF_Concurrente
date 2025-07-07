import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { ApiService } from '../../../services/api.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-form-prediccion',
  imports: [ReactiveFormsModule, CommonModule],
  templateUrl: './form-prediccion.component.html',
  styleUrl: './form-prediccion.component.css'
})
export class FormPrediccionComponent implements OnInit {
  prediccionForm!: FormGroup;
  resultado: any = null;
  enviando = false;

  constructor(
    private fb: FormBuilder,
    private api: ApiService
  ) { }

  ngOnInit(): void {
    this.prediccionForm = this.fb.group({
      hora: [null, [Validators.required]],
      departamento: ['', [Validators.required]],
      codigo_via: ['', [Validators.required]],
      kilometro: [null, [Validators.required, Validators.min(0)]],
      modalidad: ['', [Validators.required]]
    });

    this.getData();
  }

  lstValores: any;
  lstDepartamentos: any;
  lstModalidades: any;
  lstVias: any;

  getData() {
    this.api.obtenerDiccionarios().subscribe(data => {
      this.lstValores = data;
      this.lstDepartamentos = Object.entries(data.departamentos).map(
        ([key, value]) => ({ id: Number(key), nombre: value })
      );
      this.lstModalidades = Object.entries(data.modalidades).map(
        ([key, value]) => ({ id: Number(key), nombre: value })
      );
      this.lstVias = Object.entries(data.vias).map(
        ([key, value]) => ({ id: Number(key), nombre: value })
      );
    })
  }

  onSubmit(): void {
    if (this.prediccionForm.invalid) return;

    this.enviando = true;
    const req = {
      hora: Number(this.prediccionForm.get('hora')?.value),
      departamento: Number(this.prediccionForm.get('departamento')?.value),
      codigo_via: Number(this.prediccionForm.get('codigo_via')?.value),
      kilometro: Number(this.prediccionForm.get('kilometro')?.value),
      modalidad: Number(this.prediccionForm.get('modalidad')?.value)
    }
    // console.log(req);
    this.api.predecir(req).subscribe({
      next: res => {
        this.resultado = res;
        this.enviando = false;
        console.log(res);
      },
      error: err => {
        console.error(err);
        this.enviando = false;
      }
    });
  }
}
