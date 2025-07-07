import { Component, OnInit } from '@angular/core';
import { ApiService } from '../../../services/api.service';
import { ChartOptions, ChartType } from 'chart.js';
import { CommonModule } from '@angular/common';
import { NgChartsModule } from 'ng2-charts';

@Component({
  selector: 'app-estadisticas',
  imports: [CommonModule,
    NgChartsModule],
  templateUrl: './estadisticas.component.html',
  styleUrl: './estadisticas.component.css'
})
export class EstadisticasComponent  implements OnInit{
  chartDeptos: any;
  chartHoras: any;
  chartModalidad: any;

  chartOptions: ChartOptions = {
    responsive: true,
    plugins: {
      legend: {
        position: 'top'
      },
    },
  };

  constructor(private api: ApiService) {}

  ngOnInit(): void {
    this.api.obtenerEstadisticas().subscribe((data: any) => {
      // Departamentos
      const depLabels = Object.keys(data.por_departamento);
      const depValues = Object.values(data.por_departamento) as number[];
      this.chartDeptos = {
        labels: depLabels,
        datasets: [{ label: 'Accidentes por departamento', data: depValues }]
      };

      // Horas
      const horaLabels = Object.keys(data.por_hora);
      const horaValues = Object.values(data.por_hora) as number[];
      this.chartHoras = {
        labels: horaLabels,
        datasets: [{ label: 'Accidentes por hora', data: horaValues }]
      };

      // Modalidades
      const modLabels = Object.keys(data.por_modalidad);
      const modValues = Object.values(data.por_modalidad) as number[];
      this.chartModalidad = {
        labels: modLabels,
        datasets: [{ label: 'Accidentes por tipo', data: modValues }]
      };
    });
  }
}
